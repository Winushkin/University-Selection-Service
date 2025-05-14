import './ToggleSwitch.css';


// Функция работы переключателя
function ToggleSwitch({ name, checked, onChange }) {
    const handleChange = () => {
        onChange(!checked);
    };

    return (
        <label className="toggle-switch">
            <input
                type="checkbox"
                name={name}
                checked={checked}
                onChange={handleChange}
                className="toggle-input"
            />
            <span className="toggle-slider"></span>
            <span className="toggle-label">{checked ? 'Важно' : 'Не важно'}</span>
        </label>
    );
}

export default ToggleSwitch;
