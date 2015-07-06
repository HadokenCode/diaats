// Copyright 2015 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

var config struct {
	DockerHost string
	Username   string
	Password   string
	HostConfig *docker.HostConfig
	Plans      []Plan
	MongoURL   string
	DBName     string
}

type Plan struct {
	Name  string `json:"plan"`
	Image string `json:"image"`
}

func loadConfig() {
	config.DockerHost = os.Getenv("DOCKER_HOST")
	if config.DockerHost == "" {
		log.Fatal("DOCKER_HOST is required")
	}
	config.Username = os.Getenv("API_USERNAME")
	config.Password = os.Getenv("API_PASSWORD")
	hostConfigJSON := os.Getenv("DOCKER_CONFIG")
	if hostConfigJSON == "" {
		config.HostConfig = nil
	} else {
		config.HostConfig = new(docker.HostConfig)
		err := json.Unmarshal([]byte(hostConfigJSON), config.HostConfig)
		if err != nil {
			log.Fatalf("Failed to parse HOST_CONFIG: %s", err)
		}
	}
	imagePlans := os.Getenv("IMAGE_PLANS")
	if imagePlans == "" {
		log.Fatal("IMAGE_PLANS is required")
	}
	err := json.Unmarshal([]byte(imagePlans), &config.Plans)
	if err != nil {
		log.Fatalf("Failed to parse IMAGE_PLANS: %s", err)
	}
	config.MongoURL = os.Getenv("MONGODB_URL")
	if config.MongoURL == "" {
		log.Fatal("MONGODB_URL is required")
	}
	config.DBName = os.Getenv("MONGODB_DB_NAME")
	if config.DBName == "" {
		url_, err := url.Parse(config.MongoURL)
		if err != nil {
			log.Fatalf("Failed to parse MONGODB_URL: %s", err)
		}
		config.DBName = strings.TrimLeft(url_.Path, "/")
	}
}
