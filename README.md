# RNS-Babel
An utility for managing and adding languages to RNS
(Currently in the middle of development and thus not ready for use!)

### Index
* [Usage](#usage)
* [Development](#development)
  - [Testing](#testing)
* [Contributing](#contributing)
* [TODO](#todo-list)

---
## Usage
None yet :(
\
\
\
[Back](#index)

---
## Development
First make sure you have Go (Golang) installed and updated

Then clone the repo
`git clone https://github.com/Diamon0/RNS-Babel.git`

Make your changes! (or don't, I'm not one to tell you what to do)

Build
`go build .`

Profit!
(Sorry, no further documentation for now, too busy)


### Testing
Copy the necessary game files to test it (`Parser/Data/[files] and Parser/Dialog/[files]`)
Run `make test` in the main folder, feel free to add your own tests
\
\
\
[Back](#index)

---
## Contributing
Make a Fork -> ??? -> Run Tests -> Make a merge request explaining what it does, what it fixes, what it adds, and why. Or just whatever you wish to convey, I'll try to figure it out, but it helps if you explain.
\
\
\
[Back](#index)

---
## TODO List
### EXTREMELY Major
- [ ] Make the program back-up every file, probably in a folder, before making ANY change
### Major
- [ ] Finish implementing the basic operations
- [ ] Create a terminal UI (Basically a GUI but on the console)
- [X] Add support for Dialogue scripts (Diamon just took a look at it and dreaded) (Update: It wasn't as bad after I moved from a structured approach to a by-line approach)
- [ ] Create a graphical UI
- [ ] Who knows?
### Minor
- [ ] Add building/compiling to Make
\
\
\
[Back](#index)
