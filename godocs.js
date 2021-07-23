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
