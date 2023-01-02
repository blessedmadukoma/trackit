const API_URL: string = "http://localhost:8000";

const appEl = document.querySelector<HTMLElement>("#app")!;

async function getUsers(api: string): Promise<any> {
  const response = await fetch(API_URL);
  const data = await response.json();
  return data;
}

async function getUser(name: string): Promise<any> {
  const response = await fetch(API_URL);
  const data = await response.json();
  return data;
}

const user = await getUser("blessed");

appEl.innerHTML = `
<pre>
 ${JSON.stringify(user, null, 2)}
</pre>
`;

export {};
