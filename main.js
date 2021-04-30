$(document).ready(() => {
    $("body").on("dblclick", () => {
        const sel = document.getSelection();
        const content = sel.anchorNode.textContent.toString()
        const word = content.substring(sel.anchorOffset,sel.focusOffset);

        //console.log(sel);
        //console.log(content);
        //console.log(word)
        const stem = word2stem[word.toLowerCase()]
        if (stem) {
            console.log(stem)
        } else {
            console.log("word not found")
        }
        const meaning = dic[stem]
        if (meaning) {
            console.log(meaning)
        } else {
            console.log("meaning not found")
        }
    })
})
