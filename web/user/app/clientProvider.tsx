"use client";

import { ApolloClient, ApolloProvider, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
  uri: "http://localhost:8080/query",
  cache: new InMemoryCache(),
});

export type ClientProviderProps = Readonly<{
  children: React.ReactNode;
}>;

export const ClientProvider: React.FC<ClientProviderProps> = ({ children }) => {
  return <ApolloProvider client={client}>{children}</ApolloProvider>;
};
