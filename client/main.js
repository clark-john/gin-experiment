const form = document.querySelector("form")
const button = document.querySelector("form input[type=submit]")

form.onsubmit = async e => {
  const formd = new FormData(form)
  e.preventDefault();
  const response = await fetch("http://localhost:9000/submit", { 
    method: "POST",
    body: formd
  })
  console.log(await response.text());
}
