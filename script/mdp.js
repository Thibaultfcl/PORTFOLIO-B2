var mdpBg = document.getElementById('mdp');
var admin = document.getElementById('admin');
var passwordForm = document.getElementById('passwordForm');
var password = document.getElementById('password');

passwordForm.addEventListener('submit', function(e) {
    e.preventDefault();

    const correctPassword = "admin";
    const passwordValue = password.value;

    if (passwordValue === correctPassword) {
        mdpBg.style.display = 'none';
        admin.style.display = 'block';
    } else {
        alert('Mot de passe incorrect');
    }
});