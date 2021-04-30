from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response


# return (board, response)
def validate_board_id(board_id):
    if not board_id:
        return Response({
            'board_id': ErrorDetail(string='Board ID cannot be empty.',
                                    code='blank')
        }, 400)

    try:
        int(board_id)
    except ValueError:
        return Response({
            'board_id': ErrorDetail(string='Board ID must be a number.',
                                    code='invalid')
        }, 400)
