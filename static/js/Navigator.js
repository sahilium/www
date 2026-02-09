export class Navigator {
  constructor({
    containerSelector = "[data-page]",
    transitionTimeout = 250
  } = {}) {
    this.containerSelector = containerSelector;
    this.transitionTimeout = transitionTimeout;

    this._onClick = this._onClick.bind(this);
    this._onPopState = this._onPopState.bind(this);
  }

  init() {
    document.addEventListener("click", this._onClick);
    window.addEventListener("popstate", this._onPopState);
    this.onPageLoad();
  }

  get container() {
    return document.querySelector(this.containerSelector);
  }

  async navigate(url, { isPop = false } = {}) {
    if (url === location.href) return;

    const container = this.container;
    if (!container) return;

    container.classList.add("is-exiting");
    await this._waitForTransition(container);

    const html = await this._fetchPage(url);
    const { content, title } = this._parseHTML(html);

    container.innerHTML = content;
    document.title = title;

    container.classList.remove("is-exiting");
    container.classList.add("is-entering");

    requestAnimationFrame(() =>
      container.classList.remove("is-entering")
    );

    if (!isPop) {
      history.pushState({}, "", url);
    }

    this._resetScroll();
    this.onPageLoad();
  }

  /* ---------- hooks ---------- */

  onPageLoad() {
    const type = document.body.dataset.pageType;
    // extend later
  }

  /* ---------- private ---------- */

  _onClick(e) {
    const a = e.target.closest("a");
    if (!this._shouldHandleLink(e, a)) return;

    e.preventDefault();
    this.navigate(a.href);
  }

  _onPopState() {
    this.navigate(location.href, { isPop: true });
  }

  _shouldHandleLink(e, a) {
    if (!a) return false;
    if (a.origin !== location.origin) return false;
    if (a.hasAttribute("data-no-spa")) return false;
    if (e.metaKey || e.ctrlKey || e.shiftKey) return false;
    return true;
  }

  async _fetchPage(url) {
    const res = await fetch(url, {
      headers: { "X-SPA": "1" }
    });
    return res.text();
  }

  _parseHTML(html) {
    const dom = new DOMParser().parseFromString(html, "text/html");
    const next = dom.querySelector(this.containerSelector);

    return {
      content: next ? next.innerHTML : "",
      title: dom.title
    };
  }

  _waitForTransition(el) {
    return new Promise(resolve => {
      let done = false;

      const finish = () => {
        if (done) return;
        done = true;
        resolve();
      };

      el.addEventListener("transitionend", finish, { once: true });
      setTimeout(finish, this.transitionTimeout);
    });
  }

  _resetScroll() {
    document.documentElement.scrollTop = 0;
  }
}