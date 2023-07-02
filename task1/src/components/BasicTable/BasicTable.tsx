import useSignal from '../../hooks/useSignal'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import Paper from '@mui/material/Paper'
import { useEffect } from 'react'
import axios from 'axios'
import Loader from '../Loader/Loader'
import './BasicTable.css'

const API_URL = 'https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1'

interface ICurrency {
  id: string
  symbol: string
  name: string
}

export default function BasicTable() {
  const rows = useSignal<ICurrency[]>([])
  const isLoading = useSignal(false)

  useEffect(() => {
    isLoading.val = true
    axios.get(API_URL).then(({ data }) => {
      rows.val = data
      isLoading.val = false
    })
  }, [])

  if (isLoading.val) return <Loader />

  return (
    <TableContainer
      component={Paper}
      sx={{ maxWidth: 750 }}
    >
      <Table
        sx={{ width: '100%' }}
        aria-label='simple table'
      >
        <TableHead>
          <TableRow>
            <TableCell align='left'>Id</TableCell>
            <TableCell align='left'>Symbol&nbsp;(g)</TableCell>
            <TableCell align='left'>Name&nbsp;(g)</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.val.map((row, index) => (
            <TableRow
              key={row.id}
              sx={{
                '&:last-child td, &:last-child th': { border: 0 },
                backgroundColor: `${row.symbol === 'usdt' ? 'green' : index < 5 ? 'blue' : ''}`,
              }}
            >
              <TableCell
                component='th'
                scope='row'
              >
                {row.id}
              </TableCell>
              <TableCell align='left'>{row.symbol}</TableCell>
              <TableCell align='left'>{row.name}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
