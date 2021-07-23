window.onload = () => {
    console.log('window loaded');
    generateTableOfContents();
    main();
};

function main() {
  const nodesToIgnore = new Set(["PRE"]);
  const textNodes = collectTextNodes();
  for (const node of textNodes) {
    if (nodesToIgnore.has(node.parentNode.nodeName)) {
      continue;
    }
    processNode(node);
  }
}

function generateTableOfContents() {
  // get table of contents
  const container = document.querySelector(".container");
  const headers = Array.from(container.querySelectorAll("h2, h3")).filter(
    (node) => node.id !== ""
  );
  const tocItems = headers.map((node) => ({
    name: node.textContent,
    id: node.id,
    isMainHeader: node.nodeName === "H2", // H2 or H3 -> dt or dd
  }));

  // prepare DOM elements
  const table = document.createElement("table");
  table.className = "unruled";
  const tbody = document.createElement("tbody");
  table.appendChild(tbody);
  const tr = document.createElement("tr");
  tbody.appendChild(tr);
  const firstTD = document.createElement("td");
  firstTD.className = "first";
  tr.appendChild(firstTD);
  const firstDL = document.createElement("dl");
  firstTD.appendChild(firstDL);
  const secondTD = document.createElement("td");
  tr.appendChild(secondTD);
  const secondDL = document.createElement("dl");
  secondTD.appendChild(secondDL);

  // generate table of contents nodes
  const splitIndex = tocItems.length / 2 + 1;
  tocItems.forEach((tocItem, i) => {
    let headerNode;
    if (tocItem.isMainHeader) {
      headerNode = document.createElement("dt");
    } else {
      headerNode = document.createElement("dd");
      headerNode.className = "indent";
    }

    const link = document.createElement("a");
    link.href = `#${tocItem.id}`;
    link.textContent = tocItem.name;
    headerNode.appendChild(link);

    if (i < splitIndex) {
      firstDL.appendChild(headerNode);
      return;
    }
    secondDL.appendChild(headerNode);
  });

  const nav = document.getElementById("nav");
  nav.appendChild(table);
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
    const stem = word2stem[word.toLowerCase()];
    if (!stem) {
      return "";
    }
    const meaning = dic[stem];
    if (!meaning) {
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

    const trimmedWord = word.replace(/['",.:;]/g, '');
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
  // workaround to avoid making h2 > pre and h3 > pre
  if (node.parentNode.nodeName == "H2" || node.parentNode.nodeName == "H3") {
    // wrap by div
    const div = document.createElement("div");
    div.appendChild(fragment);
    node.parentNode.replaceChild(div, node);
  } else {
    node.parentNode.replaceChild(fragment, node);
  }
}
