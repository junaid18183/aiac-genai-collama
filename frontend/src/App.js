import React, { useState } from 'react';
import axios from 'axios';

function App() {
  const [inputText, setInputText] = useState('');
  const [responseText, setResponseText] = useState('');

  const handleInputChange = (event) => {
    setInputText(event.target.value);
  };

  const handleSubmit = async () => {
    try {
      const response = await axios.post(process.env.REACT_APP_BACKEND_API_URL, {
        "i_want": inputText
      });
      setResponseText(response.data);
      setInputText(''); // optional: clear the input after sending
    } catch (error) {
      // Handle errors here, if any
      console.error('An error occurred:', error);
      setResponseText('Error: Could not retrieve data.');
    }
  };
  

  return (
    <div style={{ margin: '20px' }}>
      <textarea
        value={inputText}
        onChange={handleInputChange}
        placeholder="Enter your text here"
        rows={10}
        cols={50}
      />
      <br />
      <button onClick={handleSubmit}>Send</button>
      <br />
      <textarea
        value={responseText}
        readOnly
        placeholder="Response will be shown here"
        rows={10}
        cols={50}
      />
    </div>
  );
}

export default App;
