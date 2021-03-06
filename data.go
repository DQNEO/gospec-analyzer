package main

var meaninglessTokens = map[string]bool{
	"CD":true,   // cardinal number
	"DT":true,   // determiner
	"IN":true,   // conjunction, subordinating or preposition
	"CC":true,   // conjunction, coordinating
	"PRP":true,  // pronoun, personal
	"PRP$":true, // pronoun, possessive
	"TO":true,   // infinitival to
	"WDT":true,  // wh-determiner
	"WP":true,   // wh-pronoun, personal
	"WP$":true,  // wh-pronoun, possessive
	"WRB":true,  // wh-adverb
	"MD":true,   // verb, modal auxiliary
}
