let editingContactId = null;

window.onload = function() {
    fetchContacts();
};

function fetchContacts() {
    fetch('/contacts')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.querySelector('#contactTable tbody');
            tableBody.innerHTML = '';
            data.forEach(contact => {
                const row = `<tr>
                    <td>${contact.id}</td>
                    <td>${contact.name}</td>
                    <td>${contact.phone}</td>
                    <td>${contact.email}</td>
                    <td>
                        <button onclick="editContact(${contact.id})">Edit</button>
                        <button onclick="deleteContact(${contact.id})">Delete</button>
                    </td>
                </tr>`;
                tableBody.innerHTML += row;
            });
        });
}

function addContact() {
    const name = document.getElementById('name').value;
    const phone = document.getElementById('phone').value;
    const email = document.getElementById('email').value;

    const method = editingContactId ? 'PUT' : 'POST';
    const url = editingContactId ? `/contacts/${editingContactId}` : '/contacts';

    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, phone, email }),
    }).then(response => {
        if (response.ok) {
            alert(editingContactId ? 'Contact updated' : 'Contact added');
            editingContactId = null;  // Reset after successful update
            document.getElementById('name').value = '';
            document.getElementById('phone').value = '';
            document.getElementById('email').value = '';
            fetchContacts();  // Refresh the contact list
        } else {
            alert('Error saving contact');
        }
    });
}

function editContact(id) {
    fetch(`/contacts/${id}`)
        .then(response => response.json())
        .then(contact => {
            document.getElementById('name').value = contact.name;
            document.getElementById('phone').value = contact.phone;
            document.getElementById('email').value = contact.email;
            editingContactId = contact.id;  // Set the contact ID for the update
        });
}

function deleteContact(id) {
    fetch(`/contacts/${id}`, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            alert('Contact deleted');
            fetchContacts();
        } else {
            alert('Error deleting contact');
        }
    });
}
