function delCookie() {
  const c = (document.cookie =
    "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;");
  return c;
  console.log(c);
}
