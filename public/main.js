const typed = new Typed('#code', {
  strings: [`fetch('https://daddysgotjokes.com/jokes')\n.then(res => res.json())\n.then(data => console.log(data))\n\n( ͡ᵔ ͜ʖ ͡ᵔ )`],
  typeSpeed: 50,
});

const buttons = document.querySelectorAll('#toggleButton');

buttons.forEach(function (button) {
  button.addEventListener('click', function () {
    toggleResult(this);
  });
});

function toggleResult(button) {
  const targetId = button.getAttribute('data-target');
  const resultDiv = document.getElementById(targetId);

  if (resultDiv.classList.contains('hidden')) {
    resultDiv.classList.remove('hidden');
    button.textContent = 'Hide Result';
  } else {
    resultDiv.classList.add('hidden');
    button.textContent = 'Show Result';
  }
}
