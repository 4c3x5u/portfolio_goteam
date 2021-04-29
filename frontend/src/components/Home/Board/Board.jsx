import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import { DragDropContext } from 'react-beautiful-dnd';

import AppContext from '../../../AppContext';
import ColumnsAPI from '../../../api/ColumnsAPI';
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

  const handleOnDragEnd = async (result) => {
    try {
      // Find the "source" column – the one that the task is initially in
      const source = activeBoard.columns.find((column) => (
        column.id.toString() === result.source.droppableId
      ));

      // Pop the tasks that's being moved out of the "source"
      const [item] = source.tasks.splice(result.source.index, 1);

      // Update "source" tasks' orders
      const sourceTasks = source.tasks.map((task, index) => ({
        ...task,
        order: index,
      }));

      // Find the "destination" column – the that the task is being moved into
      const destination = activeBoard.columns.find((column) => (
        column.id.toString() === result.destination.droppableId
      ));

      // Insert the task that's being moved into the "destination"
      destination.tasks.splice(result.destination.index, 0, item);

      // Update "destination" tasks' orders
      const destinationTasks = destination.tasks.map((task, index) => ({
        ...task,
        order: index,
      }));

      // Update client state (to avoid server-response wait time)
      setActiveBoard({
        ...activeBoard,
        columns: activeBoard.columns.map((column) => {
          if (column.id === destination.id) {
            return { ...column, tasks: destinationTasks };
          }
          if (column.id === source.id) {
            return { ...column, tasks: sourceTasks };
          }
          return column;
        }),
      });

      // Update the "source" and the "destination" in the database
      await ColumnsAPI.patch(source.id, sourceTasks);
      await ColumnsAPI.patch(destination.id, destinationTasks);

      // Reload the board
      await loadBoard(activeBoard.id);
    } catch (err) {
      notify(
        'Unable to update task.',
        `${err?.message || 'Server Error'}.`,
      );
      setIsLoading(false);
    }
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns
          && activeBoard.columns.map((column) => (
            column.id && (
              <Column
                key={column.id}
                id={column.id}
                order={column.order}
                tasks={column.tasks}
                handleActivate={handleActivate}
              />
            )
          ))}
      </DragDropContext>
    </Row>
  );
};

Board.propTypes = { handleActivate: PropTypes.func.isRequired };

export default Board;
