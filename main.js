$(document).ready(() => {
    $("body").on("dblclick", () => {
        const sel = document.getSelection();
        const content = sel.anchorNode.textContent.toString()
        const word = content.substring(sel.anchorOffset,sel.focusOffset);

        //console.log(sel);
        //console.log(content);
        //console.log(word)

        const found = dic[word]
        if (found) {
            console.log(found)
        } else {
            console.log("not found")
        }
    })
})
