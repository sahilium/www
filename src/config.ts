import type { SiteConfig } from "./types";

export const siteConfig: SiteConfig = {
  title: "Sahil A.",
  description: "here be dragons",
  siteUrl: "https://sahil.im",
  author: {
    name: "Sahil A.",
    bio: "I am not what happened to me. I am what I choose to become.",
    motto: "here be dragons",
  },
  nav: [
    { label: "/home", href: "/" },
    { label: "/about", href: "/about" },
    { label: "/portfolio", href: "/portfolio" },
    { label: "/blog", href: "/blog" },
    { label: "/now", href: "/now" },
    { label: "/uses", href: "/uses" },
  ],
  socials: {
    github: "https://github.com/sahilium",
    twitter: "",
    mastodon: "https://mastodon.social/@sahilium",
    matrix: "https://matrix.to/#/@nonomino:matrix.org",
    email: "mailto:hello@sahil.im",
  },
  postsPerPage: 10,
  analytics: {
    umami: {
      websiteId: "",
      src: "",
    },
  },
  rss: {
    title: "Sahil A.",
    description: "here be dragons",
  },
};
