import './App.css'
import { Box, List, ThemeIcon } from '@mantine/core'
import useSWR, { BareFetcher } from 'swr'
import AddTodo from './components/AddTodo'
import { FetcherResponse } from 'swr/dist/types'
import { CheckCircleFillIcon } from '@primer/octicons-react'

export const ENDPOINT = 'http://localhost:4000'

export interface Todo {
  id: number;
  title: string;
  body: string;
  done: Boolean;
}

const fetcher: BareFetcher<Todo[]> = async (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((r) => r.json()) as FetcherResponse<Todo[]>


function App() {
  const { data, mutate } = useSWR<Todo[]>('api/todos', fetcher)

  async function markDone(id: number) {
    const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`,{
      method: "PATCH"
    }).then(r => r.json())
    mutate(updated)
  }

  return (
    <Box
    sx={(theme) => ({
      padding: "2rem",
      width: "100%",
      maxWidth: "40rem",
      margin: "0 auto"
    })}
    >
      <List spacing="xs" size="sm" mb={12} center>
        {data?.map((todo: Todo)=> {
          return <List.Item
                    onClick={()=> markDone(todo.id)}
                    key={`todo__${todo.id}`}
                    icon={
                      todo.done? (<ThemeIcon color="teal" size={24} radius="xl">
                        <CheckCircleFillIcon size={20}></CheckCircleFillIcon>
                      </ThemeIcon>): (<ThemeIcon color="gray" size={24} radius="xl">
                        <CheckCircleFillIcon size={20}></CheckCircleFillIcon>
                      </ThemeIcon>)
                    }
                  >
                  {todo.title}
                </List.Item>
        })}
      </List>
      <AddTodo mutate={mutate}/>
    </Box>
  )
}

export default App
