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
