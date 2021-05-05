main();

function main() {
  const nodesToIgnore = new Set(["h2", "h3", "pre"]);
  const textNodes = collectTextNodes();
  for (const node of textNodes) {
    if (nodesToIgnore.has(node.parentNode.localName)) {
      continue;
    }
    processNode(node);
  }
}

function collectTextNodes() {
  const container = document.querySelector(".container")
  const r = [];
  let n;
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT, null, false);
  while ((n = walker.nextNode())) {
    r.push(n);
  }
  return r;
}

function lookupWord(word) {
    const wordLower = word.toLowerCase()
    const stem = word2stem[wordLower];
    if (!stem) {
      console.log("word '" + wordLower + " is not in word2stem")
      return "";
    }
    const meaning = dic[stem];
    if (!meaning) {
      console.log("stem '" + stem + " is not in dic")
      return "";
    }
    return meaning;
}

function processNode(node) {
  const fragment = document.createDocumentFragment();
  const words = node.textContent.split(/\s+/)
  words.forEach((word) => {
    const span = document.createElement("span");
    span.textContent = word;
    fragment.appendChild(span);
    fragment.appendChild(document.createTextNode(" "));

    const trimmedWord = word.replace(/['",.:]/g, '');
    const meaning = lookupWord(trimmedWord);
    if (!meaning) {
      return;
    }
    span.className = "word";
    const tooltip = document.createElement("div");
    tooltip.textContent = meaning;
    tooltip.className = "word-translation";
    span.appendChild(tooltip);
  });
  if (fragment.children.length > 0 && fragment.lastChild.textContent === " ") {
    fragment.removeChild(fragment.lastChild);
  }
  node.parentNode.replaceChild(fragment, node);
}
