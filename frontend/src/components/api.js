export const apiRequest = async (url, options = {}) => {
    try {
      const token = localStorage.getItem('token'); // Get token from localStorage
      const headers = {
        'Content-Type': 'application/json',
        ...(token ? { 'Authorization': `Bearer ${token}` } : {}) // Add Authorization header if token exists
      };
  
      const response = await fetch(url, {
        ...options,
        headers: {
          ...headers,
          ...(options.headers || {}) // Allow overriding headers in options
        }
      });
  
      if (!response.ok) {
        // Server returned an error
        const errorData = await response.json(); // Try to parse error response as JSON
        console.error("API Error:", errorData);
        throw new Error(`API request failed with status ${response.status}: ${JSON.stringify(errorData)}`); // Include error details in the error message
      }
  
      // Parse JSON only for successful responses
      const data = await response.json();
      return data;
  
    } catch (error) {
      console.error("API Request Error:", error);
      throw error; // Re-throw the error to be caught by the calling function
    }
  };
  