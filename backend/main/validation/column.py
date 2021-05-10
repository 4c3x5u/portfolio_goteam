from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response
import status

from .custom import CustomAPIException


# return (column, response)
def validate_column_id(column_id):
    if not column_id:
        return Response({
            'column_id': ErrorDetail(string='Column ID cannot be empty.',
                                     code='blank')
        }, 400)

    try:
        int(column_id)
    except ValueError:
        return Response({
            'column_id': ErrorDetail(string='Column ID must be a number.',
                                     code='invalid')
        }, 400)


def validate_column_id_custom(column_id):
    if not column_id:
        raise CustomAPIException('column_id',
                                 'Column ID cannot be empty.',
                                 status.HTTP_400_BAD_REQUEST)

    try:
        int(column_id)
    except ValueError:
        raise CustomAPIException('column_id',
                                 'Column ID must be a number.',
                                 status.HTTP_400_BAD_REQUEST)
