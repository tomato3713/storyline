"use client";
import { gql, useQuery } from "@apollo/client";

const GetUsers = gql`
  query GetUsers {
    users {
      id
      name
    }
  }
`;

export default function Home() {
  const { loading, error, data } = useQuery(GetUsers);

  return (
    <div>
      {loading && <p>Loading...</p>}
      {error && <p>Error: {error.message}</p>}
      {data && (
        <ul>
          <li>id: name</li>
          {data.users.map((user: { id: string; name: string }) => (
            <li key={user.id}>
              {user.id}: {user.name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
