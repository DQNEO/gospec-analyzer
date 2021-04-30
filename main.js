$(document).ready(() => {
    $("body").on("dblclick", () => {
        const sel = document.getSelection();
        const content = sel.anchorNode.textContent.toString()
        if (!content) {
            return
        }
        const word = content.substring(sel.anchorOffset,sel.focusOffset);

        //console.log(sel);
        //console.log(content);
        //console.log(word)
        const lword = word.toLowerCase()
        const stem = word2stem[lword]
        if (!stem) {
            console.log(lword + ": word not found")
            return
        }
        const meaning = dic[stem]
        if (!meaning) {
            console.log(stem + ": not defined in the dic")
            return
        }
        console.log(stem + " : " + meaning )
    })
})
