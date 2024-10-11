document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('login-form');
    const emailGroup = document.getElementById('email-group');
    const passwordGroup = document.getElementById('password-group');
    const newPasswordGroup = document.getElementById('new-password-group');
    const submitButton = document.querySelector('.login-btn');

    let currentEmail = '';

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        if (emailGroup.style.display !== 'none') {
            // Check user
            currentEmail = document.getElementById('email').value;
            try {
                const response = await fetch('http://localhost:8080/auth/check-user', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email: currentEmail }),
                });
                const data = await response.json();

                if (data.needsNewPass) {
                    emailGroup.style.display = 'none';
                    newPasswordGroup.style.display = 'block';
                    submitButton.textContent = 'Set New Password';
                } else {
                    emailGroup.style.display = 'none';
                    passwordGroup.style.display = 'block';
                    submitButton.textContent = 'Login';
                }
            } catch (error) {
                console.error('Error checking user:', error);
                alert('An error occurred. Please try again.');
            }
        } else if (newPasswordGroup.style.display !== 'none') {
            // Set new password
            const newPassword = document.getElementById('new-password').value;
            try {
                const response = await fetch('http://localhost:8080/auth/set-new-password', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email: currentEmail, new_password: newPassword }),
                });
                if (response.ok) {
                    alert('New password set successfully. Please login.');
                    newPasswordGroup.style.display = 'none';
                    passwordGroup.style.display = 'block';
                    submitButton.textContent = 'Login';
                } else {
                    alert('Failed to set new password. Please try again.');
                }
            } catch (error) {
                console.error('Error setting new password:', error);
                alert('An error occurred. Please try again.');
            }
        } else {
            // Login
            const password = document.getElementById('password').value;
            try {
                const response = await fetch('http://localhost:8080/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email: currentEmail, password: password }),
                });
                const data = await response.json();
                if (response.ok) {
                    localStorage.setItem('token', data.token);
                    alert('Login successful!');
                    // Redirect to dashboard or home page
                    // window.location.href = 'dashboard.html';
                } else {
                    alert('Login failed. Please check your credentials.');
                }
            } catch (error) {
                console.error('Error during login:', error);
                alert('An error occurred. Please try again.');
            }
        }
    });
});