package main_test

import (
	"encoding/json"
	"os/exec"
	"time"

	"github.com/concourse/time-resource/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Check", func() {
	var checkCmd *exec.Cmd

	BeforeEach(func() {
		checkCmd = exec.Command(checkPath)
	})

	Context("when executed", func() {
		var request models.CheckRequest
		var response models.CheckResponse

		BeforeEach(func() {
			request = models.CheckRequest{
				Source: models.Source{Interval: "1m"},
			}

			response = models.CheckResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := checkCmd.StdinPipe()
			Ω(err).ShouldNot(HaveOccurred())

			session, err := gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("when no version is given", func() {
			It("outputs a version containing the current time", func() {
				Ω(response).Should(HaveLen(1))
				Ω(response[0].Time.Unix()).Should(BeNumerically("~", time.Now().Unix(), 1))
			})
		})

		Context("when a version is given", func() {
			Context("with its time within the interval", func() {
				BeforeEach(func() {
					request.Version.Time = time.Now()
				})

				It("does not output any versions", func() {
					Ω(response).Should(BeEmpty())
				})
			})

			Context("with its time one interval ago", func() {
				BeforeEach(func() {
					request.Version.Time = time.Now().Add(-1 * time.Minute)
				})

				It("outputs a version containing the current time", func() {
					Ω(response).Should(HaveLen(1))
					Ω(response[0].Time.Unix()).Should(BeNumerically("~", time.Now().Unix(), 1))
				})
			})

			Context("with its time N intervals ago", func() {
				BeforeEach(func() {
					request.Version.Time = time.Now().Add(-5 * time.Minute)
				})

				It("outputs a version containing the current time", func() {
					Ω(response).Should(HaveLen(1))
					Ω(response[0].Time.Unix()).Should(BeNumerically("~", time.Now().Unix(), 1))
				})
			})
		})
	})
})