"use client";

import React, { useEffect, useState } from "react";

type Profile = {
  phone?: string;
  avatarUrl?: string;
  classLevel?: string;
  board?: string;
  preferences?: Record<string, any> | null;
};

export default function ProfilePage() {
  const [user, setUser] = useState<{ id?: string; name?: string; email?: string } | null>(null);
  const [profile, setProfile] = useState<Profile>({});
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      fetch('/api/v2/users/me').then((r) => r.json()),
      fetch('/api/v2/users/profile').then((r) => r.json()),
    ])
      .then(([me, prof]) => {
        if (!me?.error) setUser(me);
        if (!prof?.error) setProfile(prof);
      })
      .catch((err) => setMessage(String(err)))
      .finally(() => setLoading(false));
  }, []);

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setMessage(null);
    try {
      const res = await fetch('/api/v2/users/profile', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(profile),
      });
      const data = await res.json();
      if (res.ok && !data.error) {
        setMessage('Profile saved');
      } else {
        setMessage(data.error || 'Unable to save');
      }
    } catch (err) {
      setMessage(String(err));
    } finally {
      setSaving(false);
    }
  }

  function updateField<K extends keyof Profile>(key: K, value: Profile[K]) {
    setProfile((p) => ({ ...(p || {}), [key]: value }));
  }

  return (
    <main style={{ padding: 40, fontFamily: 'sans-serif' }}>
      <div style={{ maxWidth: 760, margin: '0 auto', background: '#fff', padding: 20, borderRadius: 12 }}>
        <h1>Profile</h1>
        {loading ? (
          <div>Loading…</div>
        ) : (
          <>
            <div style={{ marginBottom: 16 }}>
              <div style={{ fontSize: 14, color: '#64748b' }}>Name</div>
              <div style={{ fontWeight: 700 }}>{user?.name || '-'}</div>
              <div style={{ fontSize: 14, color: '#64748b' }}>Email</div>
              <div style={{ fontWeight: 700 }}>{user?.email || '-'}</div>
            </div>

            <form onSubmit={handleSave} style={{ display: 'grid', gap: 12 }}>
              <label>
                Phone
                <input value={profile?.phone || ''} onChange={(e) => updateField('phone', e.target.value)} style={{ width: '100%', padding: 8 }} />
              </label>

              <label>
                Avatar URL
                <input value={profile?.avatarUrl || ''} onChange={(e) => updateField('avatarUrl', e.target.value)} style={{ width: '100%', padding: 8 }} />
              </label>

              <label>
                Class / Grade
                <input value={profile?.classLevel || ''} onChange={(e) => updateField('classLevel', e.target.value)} style={{ width: '100%', padding: 8 }} />
              </label>

              <label>
                Board
                <input value={profile?.board || ''} onChange={(e) => updateField('board', e.target.value)} style={{ width: '100%', padding: 8 }} />
              </label>

              <label>
                Preferences (JSON)
                <textarea value={profile?.preferences ? JSON.stringify(profile.preferences, null, 2) : ''} onChange={(e) => {
                  try {
                    updateField('preferences', e.target.value ? JSON.parse(e.target.value) : null);
                    setMessage(null);
                  } catch (err) {
                    setMessage('Preferences JSON is invalid');
                  }
                }} rows={6} style={{ width: '100%', padding: 8 }} />
              </label>

              <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                <button type="submit" disabled={saving} style={{ background: '#2563eb', color: '#fff', padding: '8px 12px', border: 'none', borderRadius: 8 }}>
                  {saving ? 'Saving…' : 'Save profile'}
                </button>
                {message && <div style={{ color: '#334155' }}>{message}</div>}
              </div>
            </form>
          </>
        )}
      </div>
    </main>
  );
}
