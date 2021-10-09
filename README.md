# Task  | Instagram Backend API

<h3 align="center"> Instagram Backend API</h3>

  <p align="center">
    
The task is to develop a basic version of Instagram. You are only required to develop the API for the system. <br>
    <a href="https://github.com/ace1728/go_ainsta2"><strong>Explore the docs »</strong></a>
    <br>
    <br>
    <a href="https://youtu.be/8LhvL4VOChg"> Link to Task Video</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

<p align="right">(<a href="#top">back to top</a>)</p>



### Built With

* [GoLang](https://golang.org/)
* [MongoDB](https://www.mongodb.com/)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

First we will need Golang in your system and set your environment variables.
to know how to do that follow this <a href="https://www.youtube.com/watch?v=Fnnp6FBYSvo&list=PLS1QulWo1RIaRoN4vQQCYHWDuubEU8Vij&index=2">link</a>

### Prerequisites

golang version 1.0 and higher is needed

### Installation

1. Clone the repo
   
   git clone https://github.com/ace_1728/go_ainsta2.git
  
3. Install MongoDB 
  <a href="https://docs.mongodb.com/manual/installation/">Click Here :) </a>




<!-- USAGE EXAMPLES -->
## Usage
After successfully cloning the repository and installing mongodb and golang in your system. Open your terminal ,navigate to the directory having the git repository on your system and run the command <br>
``
go main.go
``
open your preferred browser and in the url box enter <br>
``
localhost:8080\users 
``
for viewing all the users in the database <br>

``
localhost:8080\users\<Id>
``
for viewing the user with User id: Id <br>
``
localhost:8080\posts 
``
for viewing all the posts <br>
``
localhost:8080\posts\<id>
``
for viewing the post with Post id: id <br>
``
localhost:8080\posts\users\<Id>
``
for viewing all the posts by the user with User id: Id <br>

### Additional Constraints/Requirements:
<pre>
The API should be developed using Go.
MongoDB should be used for storage.
Only packages/libraries <a href="https://pkg.go.dev/std">listed here</a> and <a href="https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.4.0">here</a> can be used.
</pre>
## Structures 
<h4> User </h4>
<pre>
<ol>
<li>Id</li>
<li>Name</li>
<li>Email</li>
<li>Password</li>
</ol>
<strong> The password is stored safely as it is encrypted using md5 Hash from the standard library making it <br> difficult to be reverse engineered <strong>
</pre>
<h4> Post </h4>
<pre>
<ol>
<li>User</li>
<li>Id</li>
<li>Caption</li>
<li>ImageURL</li>
<li>Posted Timestamp</li>
</ol>
</pre>

<h4>Making the server thread safe </h4>
The net/http server automatically starts a new goroutine for each client connection and executes request handlers in those goroutines. The application does not need to do anything special to serve requests concurrently.<br>

The single mux, controller and db values are used for all requests, and possibly concurrently.<br>

So i have used a mutex to protect these values in case they are not thread-safe. <br>
<pre>
type userHandlers struct {
	sync.Mutex
	store map[string]User
}
</pre>
<pre>
type postHandlers struct {
	sync.Mutex
	storep map[string]Post
}
</pre>
<!-- ROADMAP -->
## Roadmap
An HTTP JSON API capable of the following operations, <br>
-[] Create an User<br>
 &nbsp;&nbsp; -[] Should be a POST request<br>
 &nbsp;&nbsp; -[] Use JSON request body<br>
 &nbsp;&nbsp; -[] URL should be ‘/users'<br><br>
-[] Get a user using id<br>
 &nbsp;&nbsp; -[] Should be a GET request<br>
 &nbsp;&nbsp; -[] Id should be in the url parameter<br>
 &nbsp;&nbsp; -[] URL should be ‘/users/<id here>’<br><br>
-[] Create a Post<br>
 &nbsp;&nbsp; -[] Should be a POST request<br>
 &nbsp;&nbsp; -[] Use JSON request body<br>
 &nbsp;&nbsp; -[] URL should be ‘/posts'<br><br>
-[] Get a post using id<br>
 &nbsp;&nbsp; -[] Should be a GET request<br>
 &nbsp;&nbsp; -[] Id should be in the url parameter<br>
 &nbsp;&nbsp; -[] URL should be ‘/posts/<id here>’<br><br>
-[] List all posts of a user <br>
 &nbsp;&nbsp; -[] Should be a GET request<br>
 &nbsp;&nbsp; -[] URL should be ‘/posts/users/<Id here>'<br><br>

<h3> All the above listed operations have been completed and work smoothly . </h3>
<h4> <strong> The following were not completed. </strong> </h4>
<ol>
<li> Adding pagination </li>
<li> Adding test cases </li>
</ol>
<!-- CONTACT -->
## Contact
Your Name - [Khyati Chaturvedi](https://www.linkedin.com/in/khyati-chaturvedi-a7b1371b1/) - khyati.chaturvedi2019@vitstudent.ac.in<br>

Project Link: [https://github.com/ace1728/go_ainsta2](https://github.com/ace1728/go_ainsta2)<br>

