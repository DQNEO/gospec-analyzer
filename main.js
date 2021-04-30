const style = document.createElement("style");
style.innerText = `
.word {
  position: relative;
  cursor: pointer;
}
.word:hover {
  color: gray;
}
.word-label {
  display: none;
}
.word-label {
  position: absolute;
  top: -1.2rem;
  left: 0;
  color: rgba(180, 83, 9, 1);
  background-color: rgba(254, 243, 199, 1);
  font-size: 0.8rem;
  padding: 0.2rem;
  font-style: normal;
  word-break: keep-all;
  white-space: nowrap;
}
.word:hover .word-label {
  display: block;
}
`;
document.head.appendChild(style);

function textNodesUnder(el) {
  const a = [];
  let n;
  const walk = document.createTreeWalker(el, NodeFilter.SHOW_TEXT, null, false);
  while ((n = walk.nextNode())) {
    a.push(n);
  }
  return a;
}

const allTextNodes = textNodesUnder(document.querySelector(".container"));

const skippedNodes = new Set(["h2", "h3", "pre"]);

for (const node of allTextNodes) {
  const fragment = document.createDocumentFragment();
  if (skippedNodes.has(node.parentNode.localName)) {
    continue;
  }
  node.textContent.split(/\s+/).forEach((word) => {
    const span = document.createElement("span");
    span.textContent = word;
    fragment.appendChild(span);
    fragment.appendChild(document.createTextNode(" "));
    
    const trimmedWord = word.replace(/['",.]/, '');

    const lword = trimmedWord.toLowerCase();
    const stem = word2stem[lword];
    if (!stem) {
      return;
    }
    const meaning = dic[stem];
    if (!meaning) {
      return;
    }

    span.className = "word";
    const label = document.createElement("div");
    label.textContent = meaning;
    label.className = "word-label";
    span.appendChild(label);
  });
  if (fragment.children.length > 0 && fragment.lastChild.textContent === " ") {
    fragment.removeChild(fragment.lastChild);
  }
  node.parentNode.replaceChild(fragment, node);
}
