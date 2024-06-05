const fetchButton = document.querySelector('#fetchButton');
const buttons = document.querySelectorAll('#toggleButton');
const copyButtons = document.querySelectorAll('#copyButton');
const year = document.querySelector('#year');
year.innerText = new Date().getFullYear();

const typed = new Typed('#code', {
  strings: [`fetch('${window.location.href}joke')\n.then(res => res.json())\n.then(data => console.log(data))\n\n( ͡ᵔ ͜ʖ ͡ᵔ )`],
  typeSpeed: 20,
  onComplete: () => {
    fetchButton.disabled = false;
    fetchButton.innerHTML = 'Get Joke';
  },
});

fetchButton.addEventListener('click', async function () {
  fetchButton.disabled = true;
  fetchButton.innerHTML = `<span class="loader"></span>`;

  try {
    const res = await fetch(`${window.location.href}joke`);
    if (!res.ok) {
      if (res.status === 429) {
        typed.strings = [await res.text()];
        return;
      }
      throw new Error('An error occurred, please try again later.');
    }
    const data = await res.json();
    typed.strings = [JSON.stringify(data, null, 2)];
  } catch (err) {
    typed.strings = [err.message];
  } finally {
    typed.reset();
  }
});

buttons.forEach(function (button) {
  button.addEventListener('click', function () {
    toggleResult(this);
  });
});

copyButtons.forEach(function (button) {
  button.addEventListener('click', function () {
    copyContent(this);
  });
});

function toggleResult(button) {
  const targetId = button.getAttribute('data-target');
  const resultDiv = document.querySelector(`#${targetId}`);

  if (resultDiv.classList.contains('hidden')) {
    resultDiv.classList.remove('hidden');
    button.textContent = 'Hide Result';
  } else {
    resultDiv.classList.add('hidden');
    button.textContent = 'Show Result';
  }
}

function copyContent(button) {
  const svg = `<svg class="mx-auto" viewBox="173.788 213.129 92.281 73.07" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="#fff">
    <polygon
      points="262.239 213.129 258.379 213.129 254.549 213.129 250.688 213.129 250.688 216.959 246.858 216.959 246.858 220.818 242.998 220.818 242.998 224.648 239.168 224.648 239.168 228.509 235.309 228.509 235.309 232.339 231.479 232.339 231.479 236.199 227.617 236.199 227.617 240.029 223.787 240.029 223.787 243.889 219.928 243.889 219.928 247.749 216.068 247.749 216.068 251.579 212.238 251.579 212.238 255.438 208.379 255.438 208.379 259.268 204.549 259.268 204.549 263.129 200.688 263.129 200.688 259.268 196.858 259.268 196.858 255.438 192.998 255.438 192.998 251.579 189.169 251.579 189.169 247.749 185.308 247.749 185.308 243.889 181.478 243.889 177.618 243.889 173.788 243.889 173.788 247.749 173.788 251.579 173.788 255.438 173.788 259.268 173.788 263.129 177.618 263.129 177.618 266.959 181.478 266.959 181.478 270.818 185.308 270.818 185.308 274.648 189.169 274.648 189.169 278.509 192.998 278.509 192.998 282.339 196.858 282.339 196.858 286.199 200.688 286.199 204.549 286.199 208.379 286.199 208.379 282.339 212.238 282.339 212.238 278.509 216.068 278.509 216.068 274.648 219.928 274.648 219.928 270.818 223.787 270.818 223.787 266.959 227.617 266.959 227.617 263.129 231.479 263.129 231.479 259.268 235.309 259.268 235.309 255.438 239.168 255.438 239.168 251.579 242.998 251.579 242.998 247.749 246.858 247.749 246.858 243.889 250.688 243.889 250.688 240.029 254.549 240.029 254.549 236.199 258.379 236.199 258.379 232.339 262.239 232.339 262.239 228.509 266.069 228.509 266.069 224.648 266.069 220.818 266.069 216.959 266.069 213.129"
      transform="matrix(0.9999999999999999, 0, 0, 0.9999999999999999, -3.552713678800501e-15, 0)"
    />
  </svg>`;
  const targetId = button.getAttribute('data-target');
  const snippet = document.querySelector(`#${targetId}`).innerText;
  const cleanedSnippet = snippet.replace(/\s+(?=\.then)/g, '');

  navigator.clipboard.writeText(cleanedSnippet);
  button.innerHTML = svg;
  button.disabled = true;

  setTimeout(() => {
    button.innerHTML = 'Copy';
    button.disabled = false;
  }, 2500);
}

// const navIcon = document.querySelector('#nav-icon');
// navIcon.addEventListener('click', function () {
//   this.classList.toggle('open');
// });
