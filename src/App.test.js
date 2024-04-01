/* eslint-disable testing-library/prefer-screen-queries */
import { render, fireEvent, waitFor } from '@testing-library/react';
import App from './App';

test('renders JSON Validator header', () => {
  const { getByText } = render(<App />);
  const headerElement = getByText(/JSON Validator/i);
  expect(headerElement).toBeInTheDocument();
});

test('validates JSON input with Resource as *', async () => {
  const { getByText, getByRole } = render(<App />);
  
  // Mock the fetch function to simulate a server response
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ valid: false }),
    })
  );

  // Simulate typing in the textarea
  const jsonInput = `{
    "PolicyDocument": {
        "Statement": [
            {
                "Resource": "*"
            }
        ]
    }
}`;
  fireEvent.change(getByRole('textbox'), { target: { value: jsonInput } });

  // Click the button
  fireEvent.click(getByText('Check'));

  // Wait for the response to be displayed
  await waitFor(() => getByText('Response(AWS::IAM::Role Policy): {"valid":false}'));

  // Check that the response is displayed
  expect(getByText('Response(AWS::IAM::Role Policy): {"valid":false}')).toBeInTheDocument();
});

test('validates JSON input with empty Resource', async () => {
  const { getByText, getByRole } = render(<App />);
  
  // Mock the fetch function to simulate a server response
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ valid: false }),
    })
  );

  // Simulate typing in the textarea
  const jsonInput = `{
    "PolicyDocument": {
        "Statement": [
            {
                "Resource": ""
            }
        ]
    }
}`;
  fireEvent.change(getByRole('textbox'), { target: { value: jsonInput } });

  // Click the button
  fireEvent.click(getByText('Check'));

  // Wait for the response to be displayed
  await waitFor(() => getByText('Response(AWS::IAM::Role Policy): {"valid":false}'));

  // Check that the response is displayed
  expect(getByText('Response(AWS::IAM::Role Policy): {"valid":false}')).toBeInTheDocument();
});

test('validates JSON input with Resource with special signs and *', async () => {
  const { getByText, getByRole } = render(<App />);
  
  // Mock the fetch function to simulate a server response
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ valid: false }),
    })
  );

  // Simulate typing in the textarea
  const jsonInput = `{
    "PolicyDocument": {
        "Statement": [
            {
                "Resource": "res&&&*%!(&# "
            }
        ]
    }
}`;
  fireEvent.change(getByRole('textbox'), { target: { value: jsonInput } });

  // Click the button
  fireEvent.click(getByText('Check'));

  // Wait for the response to be displayed
  await waitFor(() => getByText('Response(AWS::IAM::Role Policy): {"valid":false}'));

  // Check that the response is displayed
  expect(getByText('Response(AWS::IAM::Role Policy): {"valid":false}')).toBeInTheDocument();
});

test('validates empty JSON input', async () => {
  const { getByText, getByRole } = render(<App />);
  
  // Mock the fetch function to simulate a server response
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ valid: false }),
    })
  );

  // Simulate typing in the textarea
  const jsonInput = `{}`;
  fireEvent.change(getByRole('textbox'), { target: { value: jsonInput } });

  // Click the button
  fireEvent.click(getByText('Check'));

  // Wait for the response to be displayed
  await waitFor(() => getByText('Response(AWS::IAM::Role Policy): {"valid":false}'));

  // Check that the response is displayed
  expect(getByText('Response(AWS::IAM::Role Policy): {"valid":false}')).toBeInTheDocument();
});

test('validates correct JSON input with special signs', async () => {
  const { getByText, getByRole } = render(<App />);
  
  // Mock the fetch function to simulate a server response
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ valid: true }),
    })
  );

  // Simulate typing in the textarea
  const jsonInput = `{"PolicyDocument": {
    "Statement": [
        {
            "Resource": "res&&&%!(&# "
        }
    ]
}}`;
  fireEvent.change(getByRole('textbox'), { target: { value: jsonInput } });

  // Click the button
  fireEvent.click(getByText('Check'));

  // Wait for the response to be displayed
  await waitFor(() => getByText('Response(AWS::IAM::Role Policy): {"valid":true}'));

  // Check that the response is displayed
  expect(getByText('Response(AWS::IAM::Role Policy): {"valid":true}')).toBeInTheDocument();
});