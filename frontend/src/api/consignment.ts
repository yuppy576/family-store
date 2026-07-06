import request from './request'

export interface ConsignmentData {
  itemName: string
  description?: string
  consignorId: number
  expectedPrice?: number
  commission?: number
  status: string
  isVehicle: boolean
  remark?: string
  vehicleBrand?: string
  vehicleModel?: string
  vehiclePlate?: string
  vehicleVin?: string
  vehicleYear?: number
  vehicleMileage?: number
  vehicleColor?: string
  vehicleInspection?: string
  vehicleInsurance?: string
}

// ---- Consignors ----
export function loadAllConsignors(params: any) {
  return request.get('/consignment/consignors', { params })
}
export function createConsignor(data: any) {
  return request.post('/consignment/consignors', data)
}
export function updateConsignor(id: number, data: any) {
  return request.put(`/consignment/consignors/${id}`, data)
}
export function deleteConsignor(id: number) {
  return request.delete(`/consignment/consignors/${id}`)
}

// ---- Consignment Items ----
export function loadAllConsignments(params: any) {
  return request.get('/consignment/items', { params })
}
export function createConsignment(data: any) {
  return request.post('/consignment/items', data)
}
export function updateConsignment(id: number, data: any) {
  return request.put(`/consignment/items/${id}`, data)
}
export function deleteConsignment(id: number) {
  return request.delete(`/consignment/items/${id}`)
}

// ---- Vehicles (nested under items) ----
export function createVehicle(consignmentId: number, data: any) {
  return request.post(`/consignment/items/${consignmentId}/vehicle`, data)
}
export function getVehicle(consignmentId: number) {
  return request.get(`/consignment/items/${consignmentId}/vehicle`)
}
export function updateVehicle(consignmentId: number, data: any) {
  return request.put(`/consignment/items/${consignmentId}/vehicle`, data)
}

// ---- Transfer Progress (nested under vehicles) ----
export function createTransferProgress(vehicleId: number, data: any) {
  return request.post(`/consignment/vehicles/${vehicleId}/progress`, data)
}
export function listTransferProgress(vehicleId: number) {
  return request.get(`/consignment/vehicles/${vehicleId}/progress`)
}

// ---- Settlements (nested under items) ----
export function createSettlement(consignmentId: number, data: any) {
  return request.post(`/consignment/items/${consignmentId}/settlements`, data)
}
export function listSettlements(consignmentId: number) {
  return request.get(`/consignment/items/${consignmentId}/settlements`)
}
