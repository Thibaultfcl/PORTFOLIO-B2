document.addEventListener('DOMContentLoaded', function() {
    const buttons = document.querySelectorAll('.select-project');
    const form = document.getElementById('editProjectForm');
    const projectIdInput = document.getElementById('projectId');
    const titleInput = document.getElementById('title');
    const descriptionInput = document.getElementById('description');
    const linkInput = document.getElementById('link');
    const img = document.getElementById('img');

    buttons.forEach(button => {
        button.addEventListener('click', function() {
            const id = this.getAttribute('data-id');
            const title = this.getAttribute('data-title');
            const description = this.getAttribute('data-description');
            const link = this.getAttribute('data-link');
            const imgSrc = this.getAttribute('data-picture');

            projectIdInput.value = id;
            titleInput.value = title;
            descriptionInput.value = description;
            linkInput.value = link;
            img.src = "data:image/jpeg;base64,"+imgSrc;

            form.style.display = 'block';
        });
    });
});