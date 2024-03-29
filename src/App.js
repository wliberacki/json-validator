import React, { useState } from 'react';
import axios from 'axios';

function App() {
  const [text, setText] = useState('');
  const [validationResult, setValidationResult] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await axios.post('/api/verify', { text });
    setValidationResult(response.data);
  };

  return (
    <div>
      <h1>Walidacja tekstu</h1>
      <form onSubmit={handleSubmit}>
        <input type="text" value={text} onChange={(e) => setText(e.target.value)} />
        <button type="submit">Sprawdź</button>
      </form>
      {validationResult && (
        <div>
          Wynik: {validationResult.valid ? 'Poprawny' : 'Niepoprawny'}
          {validationResult.reason && (
            <p>Przyczyna: {validationResult.reason}</p>
          )}
        </div>
      )}
    </div>
  );
};

export default App;
