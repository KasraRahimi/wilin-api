# Wilin Web API

This repository stores the source code for the REST API that serves as a backend to [my wilin website](https://www.wilin.info).

## 🏛 Introduction

Wilin is a constructed language I began working on in early 2022.
I had various goals with this language especially in regards to minimalism and phonoaesthetic appeal.
When I first began working on this language, I kept track of its lexicon in a Google Sheets.
This solution worked great at time, since the language was still malleable in its early stages and sheets offered a degree of freedom that made language development easy.
As time went on, the language's structure grew more rigid and the advantages a spreadsheet offered were being shadowed by its disadvantages.
I needed resource that would allow me to _easily_ consult the dictionary, allow me to modify it without much hassle and seemlessly add more vocabulary items.
This resource needed to be as easy to use on mobile as it was on desktop and a spreadsheet was no longer meeting my needs.

After some consideration, I decided to create a website to host my language's lexicon.
This brought many advantages with it.
For one, I can now much more easily share the dictionary with any one whose curious.
Moreover, I have much more liberty the form of the resource which means I can shape it however I like.
Finally, this website is much easier to use personally with a mobile device than the Google Sheets was.

This project exists in two separate repositories.
One repository contains the source code for the frontend of the site, that is to say the interface a user will interact with.
While this one contains the backend code, where requests will be handled and data persistance is managed.

## 💻 Tech Stack

I've made use of several technologies to help make a well-rounded API.
- GO
  - The main programming language that the API is built on.
    I chose it because of its simplicity, reliability and my familiarity with it.
- Gin
  - A backend GO framework for routing.
    GO's standard package is both powerful and simple,
    however I appreciated the abstractions Gin offered,
    especially with the way it handles route groups and middlewares.
- MySQL
  - To store user and website data.
    The data for this site is relatively simple in form,
    so I wanted to use a relational database.
    MySQL seemed like the obvious choice.
    I run MariaDB on the server that hosts the website,
    but the MySQL drivers work for it.

## ☁️ Hosting

This is a relatively small scale project, for this reason, I decided to use a very simple hosting solution
- Digital Ocean
  - They offer a variety of services, however the one I needed for my need was a simple VPS (virtual private server).
    Both the backend and frontend are served through the VPS.
- Ngnix
  - A webserver that helps in hosting my project.
    I use it to serve the frontend of my side, and it acts as a reverse proxy to redirect users to the port the server sits on.

## 🔌 Set Up

### 📜 Prerequisites

Before you can set up this project locally, there are a couple of prerequisites for it to function appropriately.
- Database
  - Make sure to have a MySQL or MariaDB server running on port 3306
  - Make sure you have a database for the server and a user with permissions to read and write to it
- Makefile (optional, but recommended)
  - I've set up a Makefile to speed up certain commands.
    You can run these commands manually on your own,
    but I highly recommend making sure Makefiles work on your device so you can benefit from it's easy of use.

### 📋 Steps
1. Clone the repository
- SSH (recommended):
```
git clone git@github.com:KasraRahimi/wilin-api.git
```
- HTTPS:
```
git clone https://github.com/KasraRahimi/wilin-api.git
```
2. Navigate to the directory
```
cd wilin-api
```
3. Configure your environment variables.
- There is a `.env.example` file at the root of this project.
  Rename it to `.env` and modify the variables to specify your database username and password.
  Make sure to also specify a secure secret key as well.
4. Install dependencies
```
go mod tidy
```
5. Run the project
- If you can run Makefiles
```
make run
```
- Otherwise you can run
```
go run .
```