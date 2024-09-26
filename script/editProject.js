document.addEventListener('DOMContentLoaded', function() {
    const buttons = document.querySelectorAll('.select-project');
    const form = document.getElementById('editProjectForm');
    const titleInput = document.getElementById('title');
    const descriptionInput = document.getElementById('description');
    const linkInput = document.getElementById('link');

    buttons.forEach(button => {
        button.addEventListener('click', function() {
            const title = this.getAttribute('data-title');
            const description = this.getAttribute('data-description');
            const link = this.getAttribute('data-link');

            titleInput.value = title;
            descriptionInput.value = description;
            linkInput.value = link;

            form.style.display = 'block';
        });
    });
});