# My Go Project

This project is a web application written in Go. It follows a clean architecture pattern and uses JWT for authentication. It's an bootcamp course microservice with a broad features to implement

## Prerequisites

Make sure you have Go installed on your machine. You can download it from the official [Go website](https://golang.org/dl/).

You will also need to install make. You can download it from the official [GNU Make website](https://www.gnu.org/software/make/).

## Installation

Follow these steps to get the project up and running:

1. Clone the repository to your local machine.

git clone https://github.com/mocolansrawung/bootcamp-auth.git
cd repo


2. Install the Go module dependencies.

go mod tidy


3. Setup your environment variables. Copy the example `.env.example` file to a new file named `.env` and replace the placeholder values with your actual values.


4. Run the application. The `make run` command will start the server.

make run


Now, you can access the web application at http://localhost:8081 (or whichever port you specified in your .env file).


Major Improvements:
1. destructure pagination query params intro struct
2. add filter based on role logged in on resolve courses (student has zero access to course list as the only one related to courses are teachers)
4. fixing hardcoded config into the use of config method