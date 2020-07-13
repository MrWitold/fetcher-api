package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MrWitold/fetcher-api/api/models"
)

const workerCount = 5

// Work struct menage task flow
type Work struct {
	jobChan     chan Job
	cancelChan  chan struct{}
	ListOfTasks map[string]*Job
}

// Job struct represents task
type Job struct {
	Interval time.Duration
	LastCall time.Time
	LinkObj  models.Link
}

func (s *Server) setupWorkers() {
	var w Work
	w.ListOfTasks = make(map[string]*Job)
	w.jobChan = make(chan Job, 10)
	s.W = &w
}

func (s *Server) initializeJobPlanner() {
	l := models.Link{}
	var err error

	listOfLinks, err := l.FindAllLinks(s.DB)
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range *listOfLinks {
		idStr := strconv.FormatUint(link.ID, 10)
		s.W.ListOfTasks[idStr] = &Job{
			Interval: time.Duration(link.Interval) * time.Second,
			LastCall: time.Now(),
			LinkObj:  link}
	}

	go s.startJobPlanner()
}

func (s *Server) startJobPlanner() {

	for {
		for _, task := range s.W.ListOfTasks {
			if time.Since(task.LastCall) > task.Interval {
				task.LastCall = time.Now()
				s.W.jobChan <- *task
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *Server) initializeWorker() {
	for i := 0; i < workerCount; i++ {
		go s.worker(s.W.jobChan, s.W.cancelChan)
	}
}

func (s *Server) removeJob(uid uint64) {
	idStr := strconv.FormatUint(uid, 10)
	delete(s.W.ListOfTasks, idStr)
}

func (s *Server) addJob(link *models.Link) {
	idStr := strconv.FormatUint(link.ID, 10)

	s.W.ListOfTasks[idStr] = &Job{
		Interval: time.Duration(link.Interval) * time.Second,
		LastCall: time.Now(),
		LinkObj:  *link}
}

func (s *Server) worker(jobChan <-chan Job, cancelChan <-chan struct{}) {
	for {
		select {
		case <-cancelChan:
			return

		case task := <-jobChan:
			go s.process(task)
		}
	}
}

func (s *Server) process(task Job) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	h := models.History{}
	resp, err := client.Get(task.LinkObj.URL)
	if err != nil {
		h.AddNewHistoryRecord(s.DB, task.LinkObj.ID, "null", 5)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	h.AddNewHistoryRecord(s.DB, task.LinkObj.ID, string(body), time.Now().Sub(start).Seconds())
}
