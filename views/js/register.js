document.addEventListener('DOMContentLoaded', function() {
  const signupForm = document.getElementById('signup-form');

  signupForm.addEventListener('submit', async function(e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const role = document.getElementById('role').value;
    const department = document.getElementById('department').value;

    try {
      const response = await fetch('http://localhost:8080/auth/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, role, department }),
      });

      if (response.ok) {
        alert('Sign up successful!');
        signupForm.reset();
      } else {
        const errorData = await response.json();
        alert(`Sign up failed: ${errorData.message}`);
      }
    } catch (error) {
      console.error('Error during sign up:', error);
      alert('An error occurred during sign up. Please try again.');
    }
  });
});
