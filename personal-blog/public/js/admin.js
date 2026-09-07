function main() {
  const elements = document.getElementsByClassName("delete");
  elements.array.forEach((e) => {
    e.addEventListener("click", (e) => {
      e.preventDefault();
      return confirm("Are you sure you want to delete this article?");
    });
  });
}

window.onload = main;
