import React, { useState, useEffect } from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom'
import { toast, ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.min.css'
import cookies from 'js-cookie'
import { jwtDecode } from 'jwt-decode'

import AppContext from './AppContext'
import InitialStates from './misc/InitialStates'
import Home from './components/Home/Home'
import Login from './components/Login/Login'
import Register from './components/Register/Register'

import 'bootstrap/dist/css/bootstrap.min.css'
import './app.sass'
import Spinner from './components/Home/Spinner/Spinner'
import TeamAPI from './api/TeamAPI'
import TasksAPI from './api/TasksAPI'
import { forEach, orderBy, some } from 'lodash'

const App = () => {
  const [isLoading, setIsLoading] = useState(false)
  const [user, setUser] = useState(InitialStates.user)
  const [team, setTeam] = useState(InitialStates.team)
  const [members, setMembers] = useState(InitialStates.members)
  const [boards, setBoards] = useState(InitialStates.boards)
  const [activeBoard, setActiveBoard] = useState(InitialStates.activeBoard)
  const notify = (header, body) => (header || body) && toast.error(
    <>
      {header && <h4>{header}</h4>}
      {body && <p>{body}</p>}
    </>,
  )

  const loadBoard = async (boardId) => {
    boardId = boardId ||
      activeBoard.id ||
      sessionStorage.getItem('board-id') ||
      ""

    let authCookie = cookies.get('auth-token')
    if (authCookie) {
      setUser(jwtDecode(authCookie))

      try {
        // if boardId is truthy, start a tasks request with it
        let tasksProm = boardId && TasksAPI.get(boardId)

        // get team - set its ID, invite token, and boards
        var teamRes = await TeamAPI.get()
        setTeam({
          id: teamRes.data.id,
          inviteToken: cookies.get("invite-token"),
        })
        // a member who isn't assigned to any board will not have any boards
        setBoards(setBoards(teamRes.data.boards ?? []))

        // it tasks request was started, await it - if not, start and await a 
        // tasks request for the first board of the team
        var tasksRes = await (
          tasksProm || TasksAPI.get(teamRes.data.boards[0].id)
        )

        let board
        if (tasksRes && tasksRes.data.length > 0) {
          // if tasks request returned any results, set the active board
          // accordingly
          board = {
            id: tasksRes.data[0].boardID,
            columns: [
              { tasks: [] }, { tasks: [] }, { tasks: [] }, { tasks: [] },
            ],
          }

          forEach(orderBy(tasksRes.data, ['order']), (task) => {
            board.columns[task.colNo].tasks.push(task)
          })
        } else {
          // if not, set the active board's board ID to the first board of team
          board = {
            id: teamRes.data.boards[0].id,
            columns: [
              { tasks: [] }, { tasks: [] }, { tasks: [] }, { tasks: [] },
            ],
          }
        }
        setActiveBoard(board)

        setMembers(teamRes.data.members.map((username) => {
          let isAdmin = username === teamRes.data.id
          let isActive = !isAdmin && board &&
            some(board.members, (m) => m == username)
          return { username, isAdmin, isActive }
        }))

      }
      catch (err) {
        // remove username if unauthorised
        if (err?.response?.status === 401) {
          setIsLoading(false)
          return
        }

        let errMsg

        if (err?.response?.data?.board) {
          notify(
            'Inactive Credentials',
            err?.response?.data?.board,
          )
          return
        }

        notify(
          'Unable to load board.',
          `${errMsg || err?.message || 'Server Error'}.`,
        )

      }
      finally { setIsLoading(false) }

    }
  }

  useEffect(() => loadBoard(), [])

  useEffect(() => (
    !cookies.get('auth-token')
    && sessionStorage.removeItem("board-id")
    && cookies.remove("invite-token")
  ), [])

  return (
    <div className="App">
      <AppContext.Provider
        value={{
          user,
          setUser,
          team,
          setTeam,
          members,
          setMembers,
          boards,
          setBoards,
          activeBoard,
          setActiveBoard,
          loadBoard,
          isLoading,
          setIsLoading,
          notify,
        }}
      >
        <Router>
          {isLoading && <Spinner />}
          <Switch>
            <Route exact path="/">
              {cookies.get('auth-token')
                ? <Home />
                : <Redirect to="/login" />}
            </Route>

            <Route path="/login">
              {!cookies.get('auth-token')
                ? <Login />
                : <Redirect to="/" />}
            </Route>

            <Route path="/register/:inviteToken?">
              {!cookies.get('auth-token')
                ? <Register />
                : <Redirect to="/" />}
            </Route>
          </Switch>
        </Router>
      </AppContext.Provider>

      <ToastContainer
        toastClassName="ErrorToast"
        position="bottom-left"
        autoClose={false}
      />
    </div>
  )
}

export default App
