import spec from '../public/openapi.json'

export default {
  paths() {
    const operationIds: string[] = []
    for (const pathItem of Object.values(spec.paths)) {
      for (const method of ['get', 'post', 'put', 'delete', 'patch'] as const) {
        const op = (pathItem as any)[method]
        if (op?.operationId) {
          operationIds.push(op.operationId)
        }
      }
    }
    return operationIds.map(operationId => ({ params: { operationId } }))
  },
}
