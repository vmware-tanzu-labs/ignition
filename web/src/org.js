async function getOrgUrl () {
  const response = await window.fetch('/api/v1/organization', {
    credentials: 'same-origin'
  })
  if (!response.ok) {
    return
  }
  const json = await response.json()
  if (!json) {
    return
  }
  return json.url
}

export { getOrgUrl }
