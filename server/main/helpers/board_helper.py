from rest_framework.serializers import ValidationError

from server.main.serializers.board.ser_board import BoardSerializer
from server.main.serializers.column.ser_column import ColumnSerializer


class BoardHelper:
    def __init__(self, name, user):
        self.name = name
        self.user = user

    def create_board(self):
        """
        Creates a board, and four columns for it.
        """
        board_serializer = BoardSerializer(
            data={'team': self.user.team_id,
                  'name': self.name}
        )
        if not board_serializer.is_valid():
            raise ValidationError({'boards': board_serializer.errors})
        board = board_serializer.save()

        board.user.add(self.user)

        # create four columns for the board
        column_data = [{'board': board.id, 'order': order}
                       for order in range(0, 4)]

        column_serializer = ColumnSerializer(data=column_data, many=True)
        if not column_serializer.is_valid():
            raise ValidationError({'columns': column_serializer.errors})

        column_serializer.save()

        return board


