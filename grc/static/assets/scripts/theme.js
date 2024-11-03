document.addEventListener("DOMContentLoaded", () => {
    const themeToggle = document.querySelector('.theme-toggle');
    if (!themeToggle) {
        console.warn('Theme toggle button not found');
        return;
    }
    
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');

    function setTheme(theme) {
        document.body.classList.toggle('dark', theme === 'dark');
        localStorage.setItem('theme', theme);
        updateThemeToggleIcon(theme);
    }

    function updateThemeToggleIcon(theme) {
        const moonPath = "M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z";
        const sunPath = "M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42";

        const svgPath = themeToggle.querySelector('path');
        if (svgPath) {
            svgPath.setAttribute('d', theme === 'dark' ? sunPath : moonPath);
        } else {
            console.warn('SVG path element not found inside theme toggle button');
        }
    }

    // Set initial theme
    const savedTheme = localStorage.getItem('theme') || (prefersDarkScheme.matches ? 'dark' : 'light');
    setTheme(savedTheme);

    themeToggle.addEventListener('click', () => {
        const currentTheme = document.body.classList.contains('dark') ? 'dark' : 'light';
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        setTheme(newTheme);
    });
});
