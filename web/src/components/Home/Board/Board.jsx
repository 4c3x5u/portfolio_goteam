import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import React, { useContext } from 'react';
import { DragDropContext } from 'react-beautiful-dnd';

import AppContext from '../../../AppContext';
import TasksAPI from '../../../api/TasksAPI';
import Column from './Column/Column';

import './board.sass';

const Board = ({ handleActivate }) => {
  const {
    activeBoard,
    loadBoard,
    setIsLoading,
    notify,
    setActiveBoard,
  } = useContext(AppContext);

  console.log("~~~ BOARD CALL")
  console.log("ACTIVE BOARD COLUMNS: " + JSON.stringify(activeBoard.columns))

  const handleOnDragEnd = async (result) => {
    if (!result?.source?.droppableId || !result?.destination?.droppableId) {
      return;
    }
    try {
      // Find the "source" column – the one that the task is initially in
      const source = activeBoard.columns.find((column) => (
        column.order.toString() === result.source.droppableId
      ));

      console.log("SOURCE: " + JSON.stringify(source))

      console.log("SOURCE TASKS: " + JSON.stringify(source.tasks))

      // Pop the tasks that's being moved out of the "source"
      const [item] = source.tasks.splice(result.source.index, 1);

      // Update "source" tasks' orders
      const sourceTasks = source.tasks.map((task, index) => ({
        ...task,
        column: parseInt(result.source.droppableId),
        order: index,
      }));

      // Find the "destination" column – the that the task is being moved into
      const destination = activeBoard.columns.find((column) => (
        column.order.toString() === result.destination.droppableId
      ));

      console.log("DESTINATION: " + JSON.stringify(destination))

      console.log("DESTINATION TASKS: " + JSON.stringify(destination.tasks))

      // Insert the task that's being moved into the "destination"
      destination.tasks.splice(result.destination.index, 0, item);

      // Update "destination" tasks' orders
      const destinationTasks = destination.tasks.map((task, index) => ({
        ...task,
        column: parseInt(result.destination.droppableId),
        order: index,
      }));

      // Update client state (to avoid server-response wait time)
      setActiveBoard({
        ...activeBoard,
        columns: activeBoard.columns.map((column, i) => {
          if (i === destination.id) {
            return { ...column, tasks: destinationTasks };
          }
          if (i === source.id) {
            return { ...column, tasks: sourceTasks };
          }
          return column;
        }),
      });

      // Update the "source" and the "destination" in the database
      sourceTasks.length > 0 && await TasksAPI.patch(sourceTasks)
      destinationTasks.length > 0 && await TasksAPI.patch(destinationTasks)
    } catch (err) {
      notify(
        'Unable to update task.',
        `${err?.message || 'Server Error'}.`,
      );
      setIsLoading(true);
      loadBoard();
    }
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns
          && activeBoard.columns.map((column, i) => (
            <Column
              id={i}
              key={i}
              order={i}
              tasks={column.tasks}
              handleActivate={handleActivate}
            />
          ))}
      </DragDropContext>
    </Row>
  );
};

Board.propTypes = { handleActivate: PropTypes.func.isRequired };

export default Board;
