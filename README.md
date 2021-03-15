# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful. 

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission. 


## Submission

1. Fork this repository and send us a link to your fork after pushing your changes. 
2. Heroku hosting - The project includes a Heroku Procfile and, in its
current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

---

## Submission notes

### Changes I made

- Fixed a bug where if a search match was close to the end (or beginning) of the TXT file, an out of bounds access was attempted (to get the previous and following 250 characters to build an excerpt)
- Initial codebase refactor so that new changes would be supported (addition of interfaces, improvement of package structure, unit tests)
- Improved result visualization (added HTML line breaks to what's sent to the frontend, highlights for searched text matches, text smoothing on the beginning and end of excerpts with '...')
- Supported case insensitive search and fuzzy search (to cope with mispellings), by using [this package](https://github.com/lithammer/fuzzysearch)
- Added Github actions pipeline to run CI on each Pull Request (worked only in this fashion), and managed to achieve a coverage of 100% in all packages other than main 

### Things I didn't have the time to do / solve in an optimal way:
- Opted for not searching the complete text passed in the search bar - I split the search upon every whitespace and run an individual search for each token.
  If I did support complete text searches, I also would have to support substring searches (based on some heuristic which doesn't come to mind right now).
  Searching by tokens seemed more straightforward (although I still have the issue of returning potentially duplicate results - maybe creating a hash for each result excerpt so I can mantain a map of those I already wrote to the response writer would work)
- Removal of stop words - this might not turn out to be an optimal user experience if implemented though. Perhaps detecting stop words and showing their related results at the bottom might be a better approach. I considered using [this dependency](https://github.com/bbalet/stopwords) but didn't quite come up with an optimal balance between ignoring those words and what happens today with the system.
- Addition of search metadata (show the number of results, search duration and whatnot to the user)
- Add integration tests for testing the REST API properly (and adding some coverage to main.go)
- Performance improvements - once I get the fuzzy search matches, I go through all the text again for each matched token (since I use a token slice for the fuzzy search, but I end up getting the actual excerpts from the original text, which is a string slice, indexes are not mappable). A solution might involve just using the token slice to get the user results (which is non-trivial, given that words are split and we can't know which whitespaces joined them originally).
