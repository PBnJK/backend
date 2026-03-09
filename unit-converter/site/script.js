const TABS = Array.from(document.getElementsByClassName("tab"));

function makeTabsInvisible() {
  TABS.forEach((tab) => {
    tab.style.display = "none";
  });
}

function switchTabs(to) {
  makeTabsInvisible();

  /* Make our desired tab visible */
  const TAB = document.getElementById(to);
  TAB.style.display = "flex";
}
