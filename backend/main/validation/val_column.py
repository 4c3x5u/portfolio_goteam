from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response


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
