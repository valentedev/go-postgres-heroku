document.getElementById("admin").addEventListener("change", adminCheckbox);
document.getElementById("ativo").addEventListener("change", ativoCheckbox);

function adminCheckbox() {
  if (this.checked) {
    this.value = "true";
  } else {
    this.value = "false";
  }
}

function ativoCheckbox() {
  if (this.checked) {
    this.value = "true";
  } else {
    this.value = "false";
  }
}
