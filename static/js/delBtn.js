// Habilita botão Deletar se email informado for igual al do DB
function BotaoDeletarUsuario(emailInput) {
  var btnSubmit = document.getElementById("deleteBtn");
  if (emailInput.value == emailDB.innerHTML.slice(8)) {
    btnSubmit.disabled = false;
  } else {
    btnSubmit.disabled = true;
  }
}
